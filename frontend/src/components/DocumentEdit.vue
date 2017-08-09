<template>
	<section>

		<form id="edit-document-form" class="row" @submit="update">

			<!-- Markdown Editor Start -->
			<div class="col-md-9">
				<h1>
					{{ heading }}
				</h1>
				<Editor></Editor>
			</div>
			<!-- Markdown Editor End -->

			<!-- Metadata Editor Start -->
			<div class="metadata-fields col-md-3">

				<FrontMatter/>
				<CommitMessageField/>

				<div class="form-group">

					<input
						type="submit"
						value="Update"
						class="btn btn-success"
						v-bind:disabled="!valid"
					/>

					<router-link :to="{name: 'document_show', params: {directory: 'documents', filename: document.filename}}" class="btn btn-text">
						Cancel
					</router-link>
				</div>

			</div>
			<!-- Metadata Editor End -->

		</form>

	</section>
</template>

<script lang="babel">
	import Editor from "../components/Editor";
	import FrontMatter from "../components/Editor/FrontMatter";
	import CommitMessageField from "../components/Editor/CommitMessageField";

	export default {
		name: "DocumentEdit",
		data() {
			return {
				markdownLoaded: false,
				valid: false,
				form: null,
			};
		},
		async created() {
			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");

			// retrieve the document and add it to vuex's store
			await this.$store.dispatch("editDocument", {directory: this.directory, filename: this.filename});
			this.markdownLoaded = true;

			this.$bus.$on("checkMetadata", () => {
				this.validate()
			});

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
			filename() {
				return this.$route.params.filename;
			},

			heading() {
				let title = this.document.title;
				if (title) {
					return title;
				} else {
					return "No title";
				}
			}
		},
		methods: {
			update(event) {
				event.preventDefault();

				this.document.update(this.commit)
					.then((response) => {
						console.debug("Document saved, redirecting to 'document_show'");
						this.redirectToShowDocument(this.directory, this.filename);
					});
			},
			redirectToShowDocument(directory, filename) {
				this.$router.push({
					name: 'document_show',
					params:{directory, filename}
				});
			},

			validate() {
				if (!this.form) {
					this.form = document.getElementById("edit-document-form");
				};
				this.valid = this.form.checkValidity();;
			}
		},
		components: {
			Editor,
			FrontMatter,
			CommitMessageField,
		}
	}
</script>