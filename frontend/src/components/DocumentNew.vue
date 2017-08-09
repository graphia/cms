<template>
	<section>

		<form id="new-document-form" class="row" @submit="create">

			<!-- Markdown Editor Start -->
			<div class="col-md-7">
				<h1>
					{{ heading }}
				</h1>
				<Editor></Editor>
			</div>
			<!-- Markdown Editor End -->

			<!-- Metadata Editor Start -->
			<div class="metadata-fields col-md-5">



				<FrontMatter/>
				<FilenameField/>
				<CommitMessageField/>

				<div class="form-group">
					<div class="btn-toolbar">

						<input
							type="submit"
							value="Create"
							class="btn btn-success"
							v-bind:disabled="!valid"
						/>

						<router-link class="btn btn-text" :to="{name: 'document_index'}">
							Cancel
						</router-link>

					</div>
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
	import FilenameField from "../components/Editor/FilenameField";

	import checkResponse from "../javascripts/response.js";

	export default {
		name: "DocumentNew",
		data() {
			return {
				valid: false,
				form: null
			};
		},
		created() {

			console.debug("new doc...");

			// initialize a fresh new document
			this.$store.dispatch("initializeDocument", this.directory);

			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");

			// when child form elements emit checkMetadata we can
			// check the validity of the form as a whole, used for
			// disabling/enabling the button
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

			heading() {
				let title = this.document.title;
				if (title) {
					return title;
				} else {
					return "New Document";
				}
			},

		},
		methods: {
			async create(event) {
				event.preventDefault();

				let response = await this.document.create(this.commit);

				if (!checkResponse(response.status)) {
					throw("could not create document");
				};

				console.debug("Document saved, redirecting to 'document_show'");
				this.redirectToShowDocument(this.document.path, this.document.filename);

			},
			redirectToShowDocument(directory, filename) {
				this.$router.push({
					name: 'document_show',
					params:{directory, filename}
				});
			},

			validate() {
				if (!this.form) {
					this.form = document.getElementById("new-document-form");
				};
				this.valid = this.form.checkValidity();;
			}

		},
		components: {
			Editor,
			CommitMessageField,
			FilenameField,
			FrontMatter
		}
	}
</script>