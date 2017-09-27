// This method taken from a gist comment by José Quintana
// https://gist.github.com/mathewbyrne/1280286#gistcomment-2005392
export default function slugifier(str) {
	const a = 'àáäâèéëêìíïîòóöôùúüûñçßÿœæŕśńṕẃǵǹḿǘẍźḧ·/_,:;'
	const b = 'aaaaeeeeiiiioooouuuuncsyoarsnpwgnmuxzh------'
	const p = new RegExp(a.split('').join('|'), 'g')

	return str.toString().toLowerCase()
		.replace(/\s+/g, '-')           // Replace spaces with -
		.replace(p, c =>
			b.charAt(a.indexOf(c)))     // Replace special chars
		.replace(/&/g, '-and-')         // Replace & with 'and'
		.replace(/[^\w\-]+/g, '')       // Remove all non-word chars
		.replace(/\-\-+/g, '-')         // Replace multiple - with single -
		.replace(/^-+/, '')             // Trim - from start of text
		.replace(/-+$/, '')             // Trim - from end of text
};