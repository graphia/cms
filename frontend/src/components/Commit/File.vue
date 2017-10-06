<template>
	<div class="card commit-file" v-bind:class="{
			'file-updated border-info': this.patch.fileUpdated(),
			'file-created border-success': this.patch.fileCreated(),
			'file-deleted border-danger': this.patch.fileDeleted()
		}"
	>

		<h2 class="card-header">
			<octicon :icon-name="this.patch.icon"/>
			<code>{{ this.path }}</code>
		</h2>

		<div class="card-body">

			<!-- file content display, different content shown depending on the git operation -->

			<Diff :patch="this.patch"/>

			<!-- end of file content display -->

		</div>

	</div>
</template>

<script lang="babel">
	import CMSPatch from '../../javascripts/models/patch.js';
	import Diff from '../Utilities/Diff';

	export default {
		name: "CommitFile",
		props: ['path', 'files'],
		data() {
			return {
				patch: null
			};
		},
		computed: {
			commitHash() {
				return this.$route.params.hash;
			}
		},
		components: {
			Diff
		},
		created() {
			this.patch = new CMSPatch(this.commitHash, this.path, this.files.old, this.files.new)
		},
	};
</script>