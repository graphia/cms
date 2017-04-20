<template>
	<section>

		<form class="row" @submit="create">

			<!-- Markdown Editor Start -->
			<div class="col-md-9">
				<h1>New Document</h1>
				<Editor></Editor>
			</div>
			<!-- Markdown Editor End -->

			<!-- Metadata Editor Start -->
			<div class="col-md-3">

				<div class="form-group">
					<label for="filename">Filename</label>
					<input name="filename" class="form-control" v-model="document.filename"/>
				</div>

				<div class="form-group">
					<label for="author">Author</label>
					<input name="author" class="form-control" v-model="document.author"/>
				</div>

				<div class="form-group">
					<label for="tags">Tags</label>
					<input name="tags" class="form-control"/>
				</div>

				<div class="form-group">
					<label for="commit-message">Commit Message</label>
					<textarea name="commit-message" class="form-control" v-model="commit.message"/>
				</div>

				<div class="form-group">
					<input type="submit" value="Update" class="btn btn-success">
				</div>

			</div>
			<!-- Metadata Editor End -->

		</form>

	</section>
</template>

<script lang="babel">
	import Editor from "../components/Editor";

	export default {
		name: "DocumentNew",
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
			}
		},
		components: {
			Editor
		}
	}
</script>