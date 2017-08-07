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
					class="form-control filename"
					type="text"
					v-model="customFilename"
			/>

			<span class="input-group-addon">
				.md
			</span>

		</div>
	</div>
</template>

<script lang="babel">
	export default {
		name: "Filename",
		data() {
			return {
				enableCustomFilename: false,
				filename: "", // filename *without* extension
			};
		},
		computed: {
			document() {
				return this.$store.state.activeDocument;
			},

			/*
			 * Deal with updates to the form's filename field depending on whether the
			 * title changes (get) or if it is modified manually (set)
			 */
			customFilename: {
				cache: true,
				get() {

					if (this.enableCustomFilename) {
						return this.filename;
					}

					let fn = "";
					if (this.document.title) {
						fn = this.slugify(this.document.title);
					}
					this.filename = fn;

					return this.filename;
				},
				set(name) {
					if (this.enableCustomFilename) {
						this.filename = this.slugify(name);
					}
				}
			}
		},
		watch: {

			/*
			 * when the filename on the form is changed (either manually or automatically)
			 * update the document's filename attribute by adding the markdown extension
			 */
			filename() {
				this.document.filename = `${this.filename}.md`;
				this.document.attachments_directory = [this.document.path, this.filename].join("/");
			}
		},
		methods: {

			// This method taken from a gist comment by José Quintana
			// https://gist.github.com/mathewbyrne/1280286#gistcomment-2005392
			slugify(text) {
				const a = 'àáäâèéëêìíïîòóöôùúüûñçßÿœæŕśńṕẃǵǹḿǘẍźḧ·/_,:;'
				const b = 'aaaaeeeeiiiioooouuuuncsyoarsnpwgnmuxzh------'
				const p = new RegExp(a.split('').join('|'), 'g')

				return text.toString().toLowerCase()
					.replace(/\s+/g, '-')           // Replace spaces with -
					.replace(p, c =>
						b.charAt(a.indexOf(c)))     // Replace special chars
					.replace(/&/g, '-and-')         // Replace & with 'and'
					.replace(/[^\w\-]+/g, '')       // Remove all non-word chars
					.replace(/\-\-+/g, '-')         // Replace multiple - with single -
					.replace(/^-+/, '')             // Trim - from start of text
					.replace(/-+$/, '')             // Trim - from end of text
			}

		}
	};
</script>

<style lang="scss" scoped>
	input.filename {
		text-align: right;
	}
</style>