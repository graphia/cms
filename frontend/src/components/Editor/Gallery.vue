<template>
	<div class="gallery">
		<h1>Gallery</h1>

		<div class="row attachments">
			<ul>
				<li v-for="(attachment, index) in document.attachments">

					<img
							class="col-md-3 img-thumbnail"
							:src="attachment.dataURI()"
							:data-size="attachment.size"
							:data-type="attachment.type"
					/>

					{{ attachment }}

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
		}
	};
</script>

<style lang="scss">
</style>
