package application

type HookMap map[string]*HookDef

func (m HookMap) ParseHooks() []error {
	errs := make([]error, 0)
	for _, v := range m {
		if err := v.parseFile(); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (m HookMap) PeekNames() []string {
	result := make([]string, 0)
	for _, v := range m {
		result = append(result, v.Name)
	}

	return result
}

func (m HookMap) FindAllMatching(acceptUrl string) []*HookDef {
	result := make([]*HookDef, 0)
	for _, v := range m {
		if v.MatchesAcceptUrl(acceptUrl) {
			result = append(result, v)
		}
	}

	return result
}
