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
		}
	};

</script>