<template>
	<div class="form-group">

		<label for="filename">Filename</label>

		<div class="input-group">

			<span class="input-group-addon">
				<label for="custom-filename" class="sr-only">Manually set the filename</label>
				<input name="custom-filename" type="checkbox" v-model="enableCustomFilename" title="Toggle custom filename"/>
			</span>

			<!-- disable tabindex when custom filename is disabled -->
			<input	:readonly="!enableCustomFilename"
					:tabindex="!enableCustomFilename ? '-1' : '0'"
					name="filename"
					class="form-control filename form-control-label"
					type="text"
					v-model="customFilename"
			/>

			<!-- only display if language isn't default, currently hardcoded to 'en' -->
			<span class="input-group-addon language-indicator" v-if="document.language != 'en'">
				.{{ document.language }}
			</span>

			<span class="input-group-addon extension-indicator">
				.md
			</span>

		</div>
	</div>
</template>

<script lang="babel">
	import Accessors from '../../Mixins/accessors';
	import slugify from '../../../javascripts/utilities/slugify.js';

	export default {
		name: "Filename",
		data() {
			return {
				enableCustomFilename: false,
				filenameBase: "", // filename *without* extension
			};
		},
		mixins: [Accessors],
		computed: {

			/*
			 * Deal with updates to the form's filename field depending on whether the
			 * title changes (get) or if it is modified manually (set)
			 */
			customFilename: {
				cache: true,
				get() {

					if (this.enableCustomFilename) {
						return this.filenameBase;
					}

					let fn = "";
					if (this.document.title) {
						fn = slugify(this.document.title);
					}
					this.filenameBase = fn;

					return this.filenameBase;
				},
				set(name) {
					if (this.enableCustomFilename) {
						this.filenameBase = slugify(name);
					}
				}
			},

			filenameWithExtension() {
				let translation = (this.document.language != "en");

				return [this.filenameBase, (translation && this.document.language), "md"]
					.filter(Boolean)
					.join(".");

			}
		},
		watch: {

			/*
			 * when the filename on the form is changed (either manually or automatically)
			 * update the document's filename attribute by adding the markdown extension, and
			 * make sure the slug matches it
			 */
			filenameBase() {
				this.document.filename = this.filenameWithExtension;
				this.document.slug = this.filenameBase;
			},

			/*
			 * if the language is changed after the title we need to trigger the updating
			 * of the filename, so the language code in the extension is present
			 */
			"this.document.language": () => {
				this.document.filename = this.filenameWithExtension;
			}
		}
	};
</script>

<style lang="scss" scoped>
	input.filename {
		text-align: right;
	}
</style>