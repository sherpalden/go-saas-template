package routes

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	adminRoutes AdminRoutes,
	userRoutes UserRoutes,
	tenantRoutes TenantRoutes,
	pcRoutes ProductCategoryRoutes,
) Routes {
	return Routes{
		adminRoutes,
		userRoutes,
		tenantRoutes,
		pcRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
