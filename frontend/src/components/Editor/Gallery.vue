<template>
	<div class="gallery">
		<h1>Gallery</h1>

		<div class="row attachments">
			<ul>
				<li v-for="(attachment, _) in document.attachments">

					<img
						class="col-md-3 img-thumbnail"
						:src="attachment.dataURI()"
						:data-size="attachment.size"
						:data-type="attachment.type"
						:data-markdown="attachment.markdownImage()"
						draggable="true"
						@dragstart="dragImage"
					/>

				</li>
			</ul>
		</div>

	</div>
</template>

<script lang="babel">
	export default {
		name: "Gallery",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		},
		methods: {
			dragImage(event) {
				console.log("dragging initiated!");
				console.debug(event);
				event.dataTransfer.setData(
					"text/plain",
					event.currentTarget.getAttribute('data-markdown')
				);
			}
		}
	};
</script>

<style lang="scss">
</style>
