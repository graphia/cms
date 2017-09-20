export default class CMSBreadcrumb {
	constructor(text, target, params) {
		this.text   = text;
		this.target = target;
		this.params = params;
		return this;
	}
};

export const HomeBreadcrumb = new CMSBreadcrumb("Home", "home", {});