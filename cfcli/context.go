package cfcli

import (
	"fmt"
	"reflect"
	"sync"

	"code.cloudfoundry.org/cli/command/common"
	"code.cloudfoundry.org/cli/plugin"
)

var context CFContext

type CFContext struct {
	cliConnection plugin.CliConnection
	cache         Cache
}

type Cache struct {
	sync.Mutex

	appList     []string
	orgList     []string
	serviceList []string
	spaceList   []string
}

func (c *Cache) Spaces() []string {
	c.Lock()
	defer c.Unlock()

	if c.spaceList == nil {
		// not initilaized
		c.spaceList = listSpaces(context.cliConnection)
	}
	spaces := make([]string, len(c.spaceList), cap(c.spaceList))
	copy(spaces, c.spaceList)
	return c.spaceList
}

func (c *Cache) Orgs() []string {
	c.Lock()
	defer c.Unlock()

	if c.orgList == nil {
		// not initilaized
		c.orgList = listOrgs(context.cliConnection)
	}
	orgs := make([]string, len(c.orgList), cap(c.orgList))
	copy(orgs, c.orgList)
	return orgs
}

func (c *Cache) Apps() []string {
	c.Lock()
	defer c.Unlock()

	if c.appList == nil {
		// not initilaized
		c.appList = listApps(context.cliConnection)
	}
	apps := make([]string, len(c.appList), cap(c.appList))
	copy(apps, c.appList)
	return apps
}

func (c *Cache) Services() []string {
	c.Lock()
	defer c.Unlock()

	if c.serviceList == nil {
		// not initilaized
		c.serviceList = listServices(context.cliConnection)
	}
	services := make([]string, len(c.serviceList), cap(c.serviceList))
	copy(services, c.serviceList)
	return services
}

func SetCFContext(cliConnection plugin.CliConnection) {
	context.cliConnection = cliConnection
	fmt.Println("Fetching visible Orgs...")
	orgs := listOrgs(cliConnection)
	fmt.Println("Fetching visible Spaces...")
	spaces := listSpaces(cliConnection)
	context.cache = Cache{spaceList: spaces, orgList: orgs}
}

func listSpaces(cliConnection plugin.CliConnection) []string {
	spaces, err := cliConnection.GetSpaces()
	if err != nil {
		return nil
	}
	spaceList := make([]string, 0, len(spaces))
	for _, space := range spaces {
		spaceList = append(spaceList, space.Name)
	}
	return spaceList
}

func listOrgs(cliConnection plugin.CliConnection) []string {
	orgs, err := cliConnection.GetOrgs()
	if err != nil {
		return nil
	}
	orgList := make([]string, 0, len(orgs))
	for _, org := range orgs {
		orgList = append(orgList, org.Name)
	}
	return orgList
}

func listApps(cliConnection plugin.CliConnection) []string {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil
	}
	appList := make([]string, 0, len(apps))
	for _, app := range apps {
		appList = append(appList, app.Name)
	}
	return appList
}

func listBuildpacks(cliConnection plugin.CliConnection) []string {
	buildpacks, err := cliConnection.GetApps()
	if err != nil {
		return nil
	}
	appList := make([]string, 0, len(buildpacks))
	for _, app := range buildpacks {
		appList = append(appList, app.Name)
	}
	return appList
}

func listServices(cliConnection plugin.CliConnection) []string {
	services, err := cliConnection.GetServices()
	if err != nil {
		return nil
	}
	serviceList := make([]string, 0, len(services))
	for _, service := range services {
		serviceList = append(serviceList, service.Name)
	}
	return serviceList
}

func listCfCommands() []string {
	t := reflect.TypeOf(common.Commands)
	commandList := make([]string, 0, t.NumField())
	for i := 1; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			commandList = append(commandList, f.Tag.Get("command"))

			alias := f.Tag.Get("alias")
			if alias != "" {
				commandList = append(commandList, alias)
			}
		}
	}
	return commandList
}
