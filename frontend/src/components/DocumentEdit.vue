<template>
	<section>

		<form class="row">

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
					<label for="tags">Tags</label>
					<input name="tags" class="form-control"/>
				</div>

				<div class="form-group">
					<label for="commit-message">Commit Message</label>
					<textarea name="commit-message" class="form-control"/>
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
			// populate $store.state.documents with docs from api

			let directory = this.$route.params.directory;
			let filename = this.$route.params.filename;

			console.debug(`retrieving document ${filename} from ${directory}`);

			this.$store.dispatch("editDocument", {directory, filename});
		},
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		}
	}
</script>