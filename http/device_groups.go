package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
)

type deviceGroupForm struct {
	Id              string
	Name            string
	Msg             string
	Devices         []*db.Device
	SelectedDevices []string
}

type deviceGroupDetails struct {
	Group *db.DeviceGroup
}

type deviceGroupsData struct {
	Count       int
	Groups      []*db.DeviceGroup
	Pagination  *Pagination
	CurrentPage int
}

func (df *deviceGroupForm) formFillIn(group *db.DeviceGroup, devices []*db.Device) {
	df.Id = group.Id
	df.Name = group.Name
	df.Devices = devices
	for _, dev := range group.Devices {
		df.SelectedDevices = append(df.SelectedDevices, dev.Id)
	}
}

func (c *HttpConfig) editDeviceGroup(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		groupErr  error
		data      = &deviceGroupForm{}
		device    = &db.Device{}
		templates = []string{deviceGroupFormTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	devsAll, err := device.GetAll(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Devices = devsAll

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.PostForm.Get("idInput")
		name := r.PostForm.Get("nameInput")
		devIds := r.PostForm["devicesInput"]

		devList := []*db.Device{}
		for _, devId := range devIds {
			d := &db.Device{}
			d.Id = devId
			devList = append(devList, d)
		}

		group := &db.DeviceGroup{
			Name:    name,
			Devices: devList,
		}
		group.Id = id

		if id == "" {
			// "id" is unset - create new group
			groupErr = group.Create(c.Db)
		} else {
			// "id" is set - update existing group
			err := group.Update(c.Db)
			if err != nil {
				data.Msg = err.Error()
			}
		}

		if groupErr != nil {
			// return data with errors if validation failed
			data.Id = id
			data.formFillIn(group, devsAll)
			data.Msg = groupErr.Error()
		} else {
			http.Redirect(w, r, "/device/groups", http.StatusFound)
			return
		}
	} else {
		// fill in the form if "id" GET parameter set
		id := r.URL.Query().Get("id")
		if id != "" {
			g := &db.DeviceGroup{}
			g.Id = id
			err = g.GetById(c.Db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(g, devsAll)
			}
		}
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) getDeviceGroups(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		g          = &db.DeviceGroup{}
		data       = &deviceGroupsData{}
		pagination = &Pagination{}
		templates  = []string{deviceGroupsTmpl, paginationTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	pageId, perPage, err := getPagionationParams(r.URL)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fetch device groups
	groupList, err := g.GetAll(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Count = len(groupList)
	if data.Count > 0 {
		chunkedDevices := chunkSliceOfObjects(groupList, perPage)
		pagination.paginate(*r.URL, pageId, len(chunkedDevices))

		if pageId-1 >= len(chunkedDevices) {
			pageId = len(chunkedDevices)
		}
		data.Pagination = pagination
		data.CurrentPage = pageId
		data.Groups = chunkedDevices[pageId-1]
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) getDeviceGroup(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		group     = &db.DeviceGroup{}
		data      = &deviceGroupDetails{}
		id        = r.URL.Query().Get("id")
		templates = []string{deviceGroupTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	group.Id = id
	err = group.GetById(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Group = group
	c.renderTemplate(w, templates, data)
}
