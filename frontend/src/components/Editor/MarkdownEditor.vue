<template>
	<div id="markdown-editor" class="form-group">

		<label for="markdown" class="sr-only">Document Contents</label>
		<textarea 	id="editor"
					name="markdown"
					class="form-control"
					rows="40"
					v-model="document.markdown"
		/>

		<Gallery/>

	</div>


</template>

<script lang="babel">

	import SimpleMDE from 'simplemde';
	import CMSFileAttachment from '../../javascripts/models/attachment.js';
	import Gallery from './Gallery';

	export default {
		name: "MarkdownEditor",
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

				let self = this;

				// Because our images might not be on the server yet, we need to
				// switch the `src` prior to displaying the preview. This function
				// makes those switches.
				let previewRender = function(text) {


					let attachments = self.document.attachments;

					let html = $.parseHTML(this.parent.markdown(text));

					console.debug("slug", self.document.slug);

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

					dropEvent.stopPropagation();
					dropEvent.preventDefault();

					console.log("Dropped!")

					// grab some information from the editor so we know where to insert
					// the image's placeholder later
					let cursor = editor.getCursor();
					let doc = editor.getDoc();
					let line = doc.getLine(cursor.line);
					let pos = {
						line: cursor.line,
						ch: line.length
					};


					// if we've dropped an image from the gallery we just care about
					// entering the placeholder into the editor
					for (var i = 0; i < dropEvent.dataTransfer.items.length; i++) {
						let item = dropEvent.dataTransfer.items[i];

						console.debug("dropped a gallery image");

						if (item.type != "text/plain") {
							console.debug(`item.type is ${item.type}, ignoring`);
							continue;
						};

						console.debug("got a text/plain, continuing")

						item.getAsString((imagePlaceholder) => {
							doc.replaceRange(`\n${imagePlaceholder}\n`, pos);
						});

					};

					// surely there's a nicer way of looping with an index in es6? 🤷
					for (var i = 0; i < dropEvent.dataTransfer.files.length; i++) {
						console.debug("dropped a filesystem image");

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

							doc.replaceRange(`\n${imagePlaceholder}\n`, pos);

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
		},
		components: {
			Gallery
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