<template>
	<section>

		<form class="row" @submit="create">

			<!-- Markdown Editor Start -->
			<div class="col-md-7">
				<h1>New Document</h1>
				<Editor></Editor>
			</div>
			<!-- Markdown Editor End -->

			<!-- Metadata Editor Start -->
			<div class="col-md-5">

				<div class="form-group">
					<label for="title">Title</label>
					<input name="title" class="form-control" v-model="document.title"/>
				</div>

				<div class="form-group">

					<!-- TODO fix this, now FrontMatter is included Title is duplicated -->
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
								v-model="customFilename"
						/>

						<span class="input-group-addon">
							.md
						</span>

					</div>
				</div>

				<FrontMatter/>

				<div class="form-group">
					<label for="commit-message">Commit Message</label>
					<textarea name="commit-message" class="form-control" v-model="commit.message"/>
				</div>

				<div class="form-group">
					<div class="btn-toolbar">
						<input type="submit" value="Update" class="btn btn-success">
						<a href="#" class="btn btn-text">Cancel</a>
					</div>
				</div>

			</div>
			<!-- Metadata Editor End -->

		</form>

	</section>
</template>

<script lang="babel">
	import Editor from "../components/Editor";
	import FrontMatter from "../components/FrontMatter";

	export default {
		name: "DocumentNew",
		data() {
			return {
				enableCustomFilename: false,
				filename: "" // filename *without* extension
			};
		},
		created() {

			console.debug("new doc...");

			// initialize a fresh new document
			this.$store.dispatch("initializeDocument", this.directory);

			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");
		},
		computed: {

			// quick access to things in the store
			document() {
				return this.$store.state.activeDocument;
			},
			commit() {
				return this.$store.state.commit;
			},

			// quick access to route params
			directory() {
				return this.$route.params.directory;
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
			create(event) {
				event.preventDefault();

				this.document.create(this.commit)
					.then((response) => {
						console.debug("Document saved, redirecting to 'document_show'");
						this.redirectToShowDocument(this.document.path, this.document.filename);
					});
			},
			redirectToShowDocument(directory, filename) {
				this.$router.push({
					name: 'document_show',
					params:{directory, filename}
				});
			},

			// This method taken from a gist comment by José Quintana
			// https://gist.github.com/mathewbyrne/1280286#gistcomment-2005392
			slugify (text) {
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
		},
		components: {
			Editor,
			FrontMatter
		}
	}
</script>

<style lang="scss" scoped>
	input.filename {
		text-align: right;
	}
</style>