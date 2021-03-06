<template>
	<div class="dropdown" v-if="anyAvailableLanguages">

		<button
			class="btn btn-primary dropdown-toggle my-2 mx-1"
			type="button"
			id="translationMenu"
			data-toggle="dropdown"
			aria-haspopup="true"
			aria-expanded="false"
		>
			Translate
		</button>

		<div class="translation-options dropdown-menu" aria-labelledby="translationMenu">

			<button @click="initiateTranslation" class="dropdown-item" v-for="(language, i) in availableLanguages" :key="i" :value="language.code">
				<span class="flag">{{ language.flag }}</span>{{ language.name }}
			</button>

		</div>

	</div>
</template>

<script lang="babel">

	import Accessors from '../Mixins/accessors';

	import CMSTranslation from '../../javascripts/models/translation.js';
	import checkResponse from "../../javascripts/response.js";
	import filenameFromLanguageCode from '../../javascripts/utilities/filename-from-language-code.js';

	export default {
		name: "Translation",
		computed: {
			languages() {
				return this.$store.state.server.translationInfo.languages;
			},
			availableLanguages() {

				return this.languages && this
					.languages
					.filter((language) => {
						return !this
							.existingTranslations
							.includes(language.code);
					});

			},
			anyAvailableLanguages() {
				return this.availableLanguages.length > 0;
			},
			existingTranslations() {
				if (this.document && this.document.translations) {
					return this.document.translations;
				}
				return [];
			}
		},
		methods: {
			async initiateTranslation(sender) {

				try {
					let code = sender.currentTarget.value;

					let filename = filenameFromLanguageCode(this.params.language_code)

					let translation = new CMSTranslation(
						this.params.directory,
						this.params.document,
						filename,
						code
					);

					let response = await translation.create();

					if (!checkResponse(response.status)) {
						throw "invalid request", response;
						return;
					};

					let json = await response.json();

					// the new filename is returned in the 'meta' field of the
					// response
					this.$store.commit("setLatestRevision", json.oid);
					this.redirectToShowDocument(this.params.directory, this.params.document, code);

					this.$store.state.broadcast.addMessage(
						"success",
						"Congratulations",
						"This is the placeholder for your translation, replace the contents " +
						"of this document with a translated version",
						10
					);
				}
				catch(err) {
					console.error("Could not create translation", err);
				}


			},
			redirectToShowDocument(directory, doc, language_code) {
				this.$router.push({
					name: 'document_show',
					params:{directory, doc, language_code}
				});
			}
		},
		mixins: [Accessors],
	};
</script>

<style lang="scss">
	button > span.flag {
		margin-right: 10px;
	};

	.dropdown {
		margin: 0rem;
	};
</style>
