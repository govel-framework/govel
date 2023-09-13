package govel

// Group Creates a new group.
func Group(prefix string, action func()) *groupModel {
	// set the variable to indicate that a group of routes exists

	newGroup := &groupModel{
		routes: make(map[string]*routeModel),
		prefix: prefix,
		parent: nil,
	}

	if inGroup {
		// edit the new group
		newGroup.parent = currentGroupConfig
		newGroup.prefix = newGroup.parent.prefix + newGroup.prefix
		newGroup.middlewares = newGroup.parent.middlewares
		newGroup.name = newGroup.parent.name

		currentGroupConfig.subGroups = append(newGroup.subGroups, newGroup)
	}

	newGroup.createGroup()
	inGroup = true

	// call the action
	action()

	if currentGroupConfig.parent == nil {
		inGroup = false
	}

	currentGroupConfig.undoGroup()

	return newGroup
}

// Name adds a name to a group.
func (gm *groupModel) Name(name string) *groupModel {
	gm.createGroup()
	gm.name = name

	// add the name to the routes
	for _, route := range gm.routes {
		route.Name(name + route.name)
	}

	// add the name to the subgroups
	for _, subGroup := range gm.subGroups {
		subGroup.Name(name)
	}

	gm.undoGroup()

	return gm
}

// Middlewares adds one or multiple middlewares to a group.
func (gm *groupModel) Middlewares(middlewares ...middlewareFunction) *groupModel {
	gm.createGroup()

	gm.middlewares = middlewares

	// add the middlewares to the routes
	for _, route := range gm.routes {
		m := route.middlewares

		m = append(m, gm.middlewares...)

		route.Middlewares(m...)
	}

	// add the middlewares to the subgroups
	for _, subGroup := range gm.subGroups {
		m := gm.middlewares

		subGroup.Middlewares(m...)
	}

	gm.undoGroup()

	return gm
}

// Internal function to "create" the group.
func (gm *groupModel) createGroup() {
	currentGroupConfig = gm
}

// Internal function to "undo" the group.
func (gm *groupModel) undoGroup() {
	currentGroupConfig = gm.parent
}
