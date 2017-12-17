<template>
	<div class="dropdown" v-if="anyAvailableLanguages">

		<button
			class="btn btn-info dropdown-toggle"
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

	export default {
		name: "Translation",
		computed: {
			languages() {
				return this.$store.state.languages;
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

					let translation = new CMSTranslation(
						this.directory,
						this.filename,
						code,
					);

					let response = await translation.create();

					if (!checkResponse(response.status)) {
						throw "invalid request", response;
						return
					}

					let json = await response.json();

					// the new filename is returned in the 'meta' field of the
					// response
					this.redirectToShowDocument(this.directory, json.meta);

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
			redirectToShowDocument(directory, filename) {
				this.$router.push({
					name: 'document_show',
					params:{directory, filename}
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
</style>
