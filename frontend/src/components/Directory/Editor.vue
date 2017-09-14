<template>
	<div id="markdown-editor" class="editor form-group">

		<label for="body" class="sr-only">Document Contents</label>
		<textarea 	id="editor"
					name="body"
					class="form-control"
					rows="40"
					v-model="directory.body"
		/>
	</div>

</template>

<script lang="babel">

	import SimpleMDE from 'simplemde';
	import CMSFileAttachment from '../../javascripts/models/attachment.js';

	export default {
		name: "MinmalMarkdownEditor", // no attachments (yet!)
		computed: {
			directory() {
				return this.$store.state.activeDirectory;
			}
		},
		mounted() {
				console.log("MarkdownEditor Created");
				this.simpleMDE = this.initializeSimpleMDE();
		},
		methods: {
			initializeSimpleMDE() {
				console.log("initializing SimpleMDE");

				let self = this;

				let simpleMDE = new SimpleMDE({
					element: document.getElementById("editor"),
					forceSync: true,
					autoFocus: true
				});

				simpleMDE.codemirror.on('change', () => {
					this.$store.state.activeDirectory.body = this.simpleMDE.value();
				});

				return simpleMDE;
			}
		},
		watch: {
			"$parent.markdownLoaded": function() {
				console.debug("syncing content");
				this.simpleMDE.value(this.directory.body);
			}
		}
	}
</script>

<style lang="scss">
	.attachments > ul > li {
		img {
			max-width: 260px;
		}
	}

	.editor-toolbar {
		background-color: #f7f7f9;
		padding: 2px 0px;

		a {
			color: black !important;
		}
	}
</style>