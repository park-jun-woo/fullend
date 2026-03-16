package crosscheck

var rules = []Rule{
	{
		Name: "OpenAPI ↔ DDL", Source: "OpenAPI", Target: "DDL",
		Requires: func(in *CrossValidateInput) bool {
			return in.OpenAPIDoc != nil && in.SymbolTable != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckOpenAPIDDL(in.OpenAPIDoc, in.SymbolTable, in.ServiceFuncs, in.SensitiveCols)
		},
	},
	{
		Name: "SSaC ↔ DDL", Source: "SSaC", Target: "DDL",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil && in.SymbolTable != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckSSaCDDL(in.ServiceFuncs, in.SymbolTable, in.DTOTypes)
		},
	},
	{
		Name: "SSaC ↔ OpenAPI", Source: "SSaC", Target: "OpenAPI",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil && in.SymbolTable != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			all := append(in.FullendPkgSpecs, in.ProjectFuncSpecs...)
			return CheckSSaCOpenAPI(in.ServiceFuncs, in.SymbolTable, in.OpenAPIDoc, all)
		},
	},
	{
		Name: "States ↔ SSaC/DDL", Source: "States", Target: "SSaC",
		Requires: func(in *CrossValidateInput) bool {
			return len(in.StateDiagrams) > 0
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckStates(in.StateDiagrams, in.ServiceFuncs, in.SymbolTable, in.OpenAPIDoc)
		},
	},
	{
		Name: "Policy ↔ SSaC/DDL/States", Source: "Policy", Target: "SSaC",
		Requires: func(in *CrossValidateInput) bool {
			return len(in.Policies) > 0
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckPolicy(in.Policies, in.ServiceFuncs, in.SymbolTable, in.StateDiagrams)
		},
	},
	{
		Name: "Scenario → OpenAPI", Source: "Scenario", Target: "OpenAPI",
		Requires: func(in *CrossValidateInput) bool {
			return len(in.HurlFiles) > 0 && in.OpenAPIDoc != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckHurlFiles(in.HurlFiles, in.OpenAPIDoc)
		},
	},
	{
		Name: "SSaC → Func", Source: "SSaC", Target: "Func",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckFuncs(in.ServiceFuncs, in.FullendPkgSpecs, in.ProjectFuncSpecs, in.SymbolTable, in.OpenAPIDoc)
		},
	},
	{
		Name: "Config → OpenAPI", Source: "Config", Target: "OpenAPI",
		Requires: func(in *CrossValidateInput) bool {
			return in.OpenAPIDoc != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckMiddleware(in.Middleware, in.OpenAPIDoc)
		},
	},
	{
		Name: "SSaC → Config", Source: "SSaC", Target: "Config",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckClaims(in.ServiceFuncs, in.Claims)
		},
	},
	{
		Name: "Policy → Config (claims)", Source: "Policy", Target: "Config",
		Requires: func(in *CrossValidateInput) bool {
			return len(in.Policies) > 0 && in.Claims != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckClaimsRego(in.Policies, in.Claims)
		},
	},
	{
		Name: "DDL → SSaC (coverage)", Source: "DDL", Target: "SSaC",
		Requires: func(in *CrossValidateInput) bool {
			return in.SymbolTable != nil && in.ServiceFuncs != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckDDLCoverage(in.SymbolTable, in.ServiceFuncs, in.Archived)
		},
	},
	{
		Name: "SSaC Queue", Source: "SSaC", Target: "",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckQueue(in.ServiceFuncs, in.QueueBackend)
		},
	},
	{
		Name: "SSaC → Authz", Source: "SSaC", Target: "Func",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckAuthz(in.ServiceFuncs, in.AuthzPackage)
		},
	},
	{
		Name: "DDL Sensitive", Source: "DDL", Target: "",
		Requires: func(in *CrossValidateInput) bool {
			return in.SymbolTable != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckSensitiveColumns(in.SymbolTable, in.SensitiveCols, in.NoSensitiveCols)
		},
	},
	{
		Name: "Func → SSaC (coverage)", Source: "Func", Target: "SSaC",
		Requires: func(in *CrossValidateInput) bool {
			return in.ServiceFuncs != nil && len(in.ProjectFuncSpecs) > 0
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckFuncCoverage(in.ServiceFuncs, in.ProjectFuncSpecs)
		},
	},
	{
		Name: "Policy → Config (roles)", Source: "Policy", Target: "Config",
		Requires: func(in *CrossValidateInput) bool {
			return len(in.Policies) > 0 && len(in.Roles) > 0
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckRoles(in.Policies, in.Roles)
		},
	},
	{
		Name: "OpenAPI Constraints", Source: "OpenAPI", Target: "DDL",
		Requires: func(in *CrossValidateInput) bool {
			return in.SymbolTable != nil && in.SymbolTable.RequestSchemas != nil && in.ServiceFuncs != nil
		},
		Check: func(in *CrossValidateInput) []CrossError {
			return CheckOpenAPIConstraints(in)
		},
	},
}
