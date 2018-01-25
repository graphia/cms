<template>
	<div>
		<!-- file content display, different content shown depending on the git operation -->

		<div
			:id="hashPresent ? `diff-${hash}` : null"
			class="raw"
			:class="{
				'collapse': collapse,
				'multi-collapse': collapse,
				'file-updated': patch.fileUpdated(),
				'file-created': patch.fileCreated(),
				'file-deleted': patch.fileDeleted()
			}"
		>

			<div v-if="isImage(patch)">
				<div v-if="patch.fileUpdated()" class="row diff" :src="imageSrc(patch.newFile, patch.filename)">
					<!-- image has been updated, show old and new side by side -->
					<div class="col rounded-left bg-danger col-md-6 patch-image">
						<img class="img-fluid" :src="imageSrc(patch.oldFile, patch.filename)"/>
					</div>

					<div class="col rounded-right bg-success col-md-6 patch-image">
						<img class="img-fluid" :src="imageSrc(patch.newFile, patch.filename)"/>
					</div>

				</div>
				<div v-else-if="patch.fileCreated()" class="col rounded bg-success col-md-12 patch-image">
				 	<img class="img-fluid" :src="imageSrc(patch.newFile, patch.filename)"/>
				</div>
				<div v-else-if="patch.fileDeleted()" class="col rounded bg-danger col-md-12 patch-image">
				 	<img class="img-fluid" :src="imageSrc(patch.oldFile, patch.filename)"/>
				</div>

			</div>
			<div v-else>

				<pre v-if="patch.fileUpdated()" class="diff" v-html="patch.diff()"/>
				<pre v-else-if="patch.fileCreated()" class="diff">{{ patch.newFile }}</pre>
				<pre v-else-if="patch.fileDeleted()" class="diff">{{ patch.oldFile }}</pre>

			</div>

		</div>
		<!-- end of file content display -->
	</div>
</template>

<script lang="babel">

	export default {
		name: "Diff",
		props: ["patch", "hash", "collapsible"],
		created() {

			if (!this.patch) {
				console.warn("A patch is required to display a diff!");
			};

			if (this.collapsible == undefined) {
				this.collapse = false;
			} else {
				this.collapse = this.collapsible;
			};

		},
		computed: {
			hashPresent() {
				return !!this.hash;
			}
		},
		methods: {
			isImage(patch) {
				let extensions = this.$config.image_extensions;
				return extensions.some((ext) => {return patch.filename.endsWith(ext)})
			},

			imageSrc(str, fn) {
				return `data:image/png;base64,${str}`
			}
		}
	};

</script>

<style lang="scss">
.patch-image {
	display: flex;
	justify-content: center;
	align-items: center;
}
</style>