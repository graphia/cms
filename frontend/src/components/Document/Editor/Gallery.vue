<template>

	<div class="card-body image-list row">
		<div class="col-xs-6" v-for="(attachment, _) in document.attachments">

			<img
				:src="attachment.dataURI()"
				:data-size="attachment.size"
				:data-type="attachment.type"
				:data-markdown="attachment.markdownImage()"
				draggable="true"
				@dragstart="dragImage"
			/>

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
	.image-list {
		overflow-x: auto;
	}
</style>
