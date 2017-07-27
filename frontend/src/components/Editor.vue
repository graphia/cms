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

					<img
						class="col-md-3 img-thumbnail"
						:src="attachment.dataURI()"
						:data-size="attachment.size"
						:data-type="attachment.type"
					/>

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


				simpleMDE.codemirror.on('drop', (editor, dropEvent) => {
					console.log("dropped!");

					let cursor = editor.getCursor();
					let doc = editor.getDoc();
					let line = doc.getLine(cursor.line);
					let pos = { // create a new object to avoid mutation of the original selection
						line: cursor.line,
						ch: line.length // set the character position to the end of the line
					};


					console.log(cursor);
					dropEvent.stopPropagation();
					dropEvent.preventDefault();

					for (var item in dropEvent.dataTransfer.items) {
						console.log("item:", item);
					};

					// surely there's a nicer way of looping with an index in es6? 🤷
					for (var i = 0; i < dropEvent.dataTransfer.files.length; i++) {

						let file = dropEvent.dataTransfer.files[i];
						let reader = new FileReader();

						reader.onloadend = (onloadendEvent) => {

							// add a CMSFileAtachment to the attachments list
							let attachment = new CMSFileAttachment(
								file,
								onloadendEvent.target.result,
								{base64Encoded: true}
							);

							this.document.addAttachment(attachment);

							let image_placeholder = attachment.markdownImage();

							doc.replaceRange(`\n${image_placeholder}\n`, pos);

						};

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

<style lang="scss">
.attachments > ul > li {
	img {
		max-width: 260px;
	}
}
</style>