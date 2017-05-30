<template>
	<section>

		<form class="row" @submit="update">

			<!-- Markdown Editor Start -->
			<div class="col-md-9">
				<h1>Editing {{ document.filename }}</h1>
				<Editor></Editor>
			</div>
			<!-- Markdown Editor End -->

			<!-- Metadata Editor Start -->
			<div class="col-md-3">

				<FrontMatter/>

				<div class="form-group">
					<label for="commit-message">Commit Message</label>
					<textarea name="commit-message" class="form-control" v-model="commit.message"/>
				</div>

				<div class="form-group">
					<input type="submit" value="Update" class="btn btn-success">

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
	import FrontMatter from "../components/FrontMatter";

	export default {
		name: "DocumentEdit",
		data() {
			return {markdownLoaded: false};
		},
		async created() {
			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");

			// retrieve the document and add it to vuex's store
			await this.$store.dispatch("editDocument", {directory: this.directory, filename: this.filename});
			this.markdownLoaded = true;

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
			}
		},
		components: {
			Editor,
			FrontMatter
		}
	}
</script>