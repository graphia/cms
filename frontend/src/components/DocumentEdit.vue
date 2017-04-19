<template>
	<section>

		<form class="row" @submit="update">

			<!-- Markdown Editor Start -->
			<div class="col-md-9">
				<h1>Editing {{ document.filename }}</h1>

				<div class="form-group">
					<label for="markdown" class="sr-only">Document Contents</label>
					<textarea name="markdown" class="form-control" rows="40" v-model="document.markdown"/>
				</div>
			</div>
			<!-- Markdown Editor End -->

			<!-- Metadata Editor Start -->
			<div class="col-md-3">

				<div class="form-group">
					<label for="title">Title</label>
					<input name="title" class="form-control" v-model="document.title"/>
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
	export default {
		name: "DocumentEdit",
		created() {
			// retrieve the document and add it to vuex's store
			this.$store.dispatch("editDocument", {directory: this.directory, filename: this.filename});

			// set up a fresh new commit
			this.$store.dispatch("initializeCommit")
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
				console.debug('Update clicked!');

				this.document.update(this.commit);
			},
			redirectToShowDocument() {
				this.$router.push({
					name: 'document_show',
					params:{
						filename: this.filename,
						directory: this.directory
					}
				});
			}
		}
	}
</script>