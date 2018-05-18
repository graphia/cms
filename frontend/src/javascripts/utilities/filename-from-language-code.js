import store from '../store.js';

export default function filename(language_code, filename="index", ext="md") {

	const dl = store.state.server.translationInfo.defaultLanguage;

	if (language_code && language_code != dl) {
		return `${filename}.${language_code}.${ext}`;
	};

	return `${filename}.${ext}`;

}
