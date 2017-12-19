<template>
	<div id="markdown-editor" class="form-group">

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
	import CMSFileAttachment from '../../../javascripts/models/attachment.js';
	import Accessors from '../../Mixins/accessors';

	export default {
		name: "MarkdownEditor",
		mixins: [Accessors],
		data() {
			return {
				count: 0
			}
		},
		mounted() {
				this.simpleMDE = this.initializeSimpleMDE();
		},
		methods: {
			initializeSimpleMDE() {

				let self = this;

				// Because our images might not be on the server yet, we need to
				// switch the `src` prior to displaying the preview. This function
				// makes those switches.
				let previewRender = function(text) {

					let attachments = self.document.attachments;

					let html = $.parseHTML(this.parent.markdown(text));

					$(html)
						.find('img')
						.each(function(_, element) {

							if ($(element)
								.attr('src')
								.startsWith("images/")) {

								let attachment = attachments
									.find(
										function(a) {
											return a.relativePath() === $(element).attr('src')
										}
									);

								if (!attachment) {
									console.warn("No attachment found matching", element);
								};

								$(element).attr('src', attachment.dataURI());

							};
						});

					return html
						.map((e) => {return e.outerHTML})
						.join("");

				}

				let simpleMDE = new SimpleMDE({
					element: document.getElementById("editor"),
					forceSync: true,
					autoFocus: true,
					dragDrop: true,
					allowDropFileTypes: ["image/jpeg", "image/jpg", "image/png", "image/gif"],
					previewRender
				});

				simpleMDE.codemirror.on('change', () => {
					this.$store.state.activeDocument.markdown = this.simpleMDE.value();
				});

				simpleMDE.codemirror.on('drop', async (editor, dropEvent) => {


					this.count++;


					for (let type of event.dataTransfer.types) {
						console.debug("type:", type);
					};

					/*
					 * Draging text externally (dragging text from another file): types
					 * has "text/plain" and "text/html"
					 * Draging text internally (dragging text to another line): types
					 * has just "text/plain"
					 * Draging a file: types has "Files"
					 * Draging a url: types has "text/plain" and "text/uri-list"
					 */

					// grab some information from the editor so we know where to insert
					// the image's placeholder later
					let cursor = editor.getCursor();
					let doc = editor.getDoc();
					let line = doc.getLine(cursor.line);
					let pos = {
						line: cursor.line,
						ch: line.length
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

							let imagePlaceholder = attachment.markdownImage();

							doc.replaceRange(`${imagePlaceholder}\n`, pos);

						};

						reader.readAsDataURL(file);

					};

				});

				return simpleMDE;
			}
		},
		watch: {
			"$parent.markdownLoaded": function() {
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

	.editor-toolbar {
		background-color: #f7f7f9;
		padding: 2px 0px;

		a {
			color: black !important;
		}
	}
</style>