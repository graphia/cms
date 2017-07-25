<template>
	<div id="editor-container" class="form-group">

		<label for="markdown" class="sr-only">Document Contents</label>
		<textarea 	id="editor"
					name="markdown"
					class="form-control"
					rows="40"
					v-model="document.markdown"
		/>

		<div class="row attachments">
			<ul>
				<li v-for="(attachment, index) in document.attachments">
					<h2>{{ attachment.name }}</h2>

					<img :src="attachment.dataURI()"/>

				</li>

			</ul>
		</div>

	</div>


</template>

<script lang="babel">

	import SimpleMDE from 'simplemde';
	import CMSFileAttachment from '../javascripts/models/attachment.js';

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
					autoFocus: true,
					dragDrop: true,
					allowDropFileTypes: true // ["image/jpeg", "image/jpg", "image/png", "image/gif"]
				});

				simpleMDE.codemirror.on('change', () => {
					this.$store.state.activeDocument.markdown = this.simpleMDE.value();
				});


				simpleMDE.codemirror.setOption('onDragEvent', function(editor, event) {
					if (event.type === "drop") {
						event.stopPropagation();
						event.preventDefault();
					};
				});


				simpleMDE.codemirror.on('drop', (editor, event) => {
					console.log("dropped!");

					event.stopPropagation();
					event.preventDefault();

					for (var item in event.dataTransfer.items) {
						console.log("item:", item);
					}

					// surely there's a nicer way of looping with an index in es6? 🤷
					for (var i = 0; i < event.dataTransfer.files.length; i++) {

						let file = event.dataTransfer.files[i];
						let reader = new FileReader();

						reader.onloadend = (event) => {

							// add a CMSFileAtachment to the attachments list
							return this.document.addAttachment(
								new CMSFileAttachment(
									file,
									event.target.result
								)
							);
						}

						reader.readAsDataURL(file);

					};

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