<template>

	<div class="card-body image-list row">
		<div class="col-xs-3" v-for="(attachment, i) in document.attachments" :key="i">

			<img
				:src="attachment.dataURI()"
				:data-size="attachment.size"
				:data-type="attachment.type"
				:data-markdown="attachment.markdownImage()"
				class="img-thumbnail rounded"
				draggable="true"
				height="100px"
				@dragstart="dragImage"
			/>

		</div>
	</div>

</template>

<script lang="babel">
	import Accessors from '../../Mixins/accessors';

	export default {
		name: "Gallery",
		mixins: [Accessors],
		methods: {
			dragImage(event) {
				event.dataTransfer.clearData();

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
