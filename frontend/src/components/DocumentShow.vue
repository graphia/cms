<template>
	<section class="row">

		<article id="document-content" class="col-md-8">
			<div class="content" v-html="relativeHTML"/>
		</article>

		<aside class="col-md-4">
			<div class="card document-metadata">
				<div class="card-block">
					<dl>

						<dt>Title</dt>
						<dd>{{ document.title  }}</dd>


						<dt>Description</dt>
						<dd>{{ document.synopsis }}</dd>

						<dt>Version</dt>
						<dd>{{ document.version }}</dd>

						<dt>Tags</dt>
						<dd>
							<span v-for="tag in document.tags" class="tag badge badge-primary">
								{{ tag }}
							</span>
						</dd>
					</dl>

					<div class="btn-toolbar" role="toolbar">
						<router-link class="btn btn-success mr-2" :to="{name: 'document_edit', params: {directory: 'documents', filename: this.filename}}">
							Edit
						</router-link>

						<button type="button" @click="destroy" class="btn btn-danger mr-2">
							Delete
						</button>
					</div>
				</div>
			</div>
		</aside>
	</section>
</template>

<style scoped lang="scss">
	aside {
		margin: 2em 0em;
	}

	.document-metadata {
		span.tag {
			margin-right: 0.6em;
		}
	}
</style>

<script lang="babel">
	export default {
		name: "DocumentShow",
		created() {
			// populate $store.state.documents with docs from api

			let directory = this.directory;
			let filename = this.filename;

			console.debug(`retrieving document ${filename} from ${directory}`);

			this.$store.dispatch("getDocument", {directory, filename});

			// create a commit to be populated/used if delete is clicked
			this.$store.dispatch("initializeCommit")

		},
		computed: {
			directory() {
				return this.$route.params.directory;
			},
			filename() {
				return this.$route.params.filename;
			},
			document() {
				return this.$store.state.activeDocument;
			},
			commit() {
				return this.$store.state.commit;
			},
			// Amend any relative links or images to point at the
			// correct resource
			relativeHTML() {

				let attachmentsDir = this.document.slug;
				let html = $.parseHTML(this.document.html);

				$(html)
					.find('img')
					.each(function(_, image) {
						if ($(image)
							.attr('src')
							.startsWith("images")) {
								let src = $(image).attr('src');
								$(image).attr('src', [attachmentsDir, src].join("/"));
						};
					});

				return html
						.map((e) => {return e.outerHTML})
						.join("");


			}
		},
		methods: {
			destroy(event) {
				event.preventDefault();
				console.debug("delete clicked!");

				let commit = this.commit;
				let file = this.document;

				this.document.destroy(this.commit)
					.then((response) => {
						console.debug("File deleted, redirecting to document index");
						this.redirectToDirectoryIndex(this.directory);
					});

			},
			redirectToDirectoryIndex(directory) {
				this.$router.push({
					name: 'document_index',
					params:{directory}
				});
			}
		}
	}
</script>