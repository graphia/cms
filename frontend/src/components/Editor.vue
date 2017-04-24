<template>
	<div id="editor-container" class="form-group">

		<label for="markdown" class="sr-only">Document Contents</label>
		<textarea 	id="editor"
					name="markdown"
					class="form-control"
					rows="40"
					v-model="document.markdown"
		/>
	</div>
</template>

<script lang="babel">

	import SimpleMDE from 'simplemde';

	export default {
		name: "Editor",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		},
		mounted() {
				console.log("MarkdownEditor Created");
				this.simpleMDE = this.initializeSimpleMDE();
		},
		methods: {
			initializeSimpleMDE() {
				console.log("initializing SimpleMDE");

				let simpleMDE = new SimpleMDE({
					element: document.getElementById("editor"),
					forceSync: true,
					autoFocus: true
				});

				simpleMDE.codemirror.on('change', () => {
					this.$store.state.activeDocument.markdown = this.simpleMDE.value();
				});

				return simpleMDE;
			}
		},
		watch: {
			"$parent.markdownLoaded": function() {
				console.debug("syncing content");
				this.simpleMDE.value(this.document.markdown);
			}
		}
	}
</script>