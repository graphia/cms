<template>
	<div class="card file" v-bind:class="{
			'file-updated card-outline-info': fileUpdated,
			'file-created card-outline-success': fileCreated,
			'file-deleted card-outline-danger': fileDeleted
		}"
	>

		<h2 class="card-header">
			<code>{{ this.path }}</code>
		</h2>

		<div class="card-block">

			<!-- file content display, different content shown depending on the git operation -->

			<div v-if="this.fileUpdated">
				<h3 class="card-subtitle">File updated</h3>
				<pre v-html="this.diff"/>
			</div>

			<div v-else-if="this.fileCreated" class="file-created">
				<h3 class="card-subtitle">File Created</h3>
				<pre>{{ this.newFile }}</pre>
			</div>

			<div v-else-if="this.fileDeleted" class="file-deleted">
				<h3 class="card-subtitle">File Deleted</h3>
				<pre>{{ this.oldFile }}</pre>
			</div>

			<!-- end of file content display -->

		</div>

	</div>
</template>

<script lang="babel">
	import Diff from 'text-diff';
	export default {
		name: "CommitFile",
		props: ['path', 'files'],
		data() {
			return {
				oldFile: null,
				oldFilePresent: null,
				newFile: null,
				newFilePresent: null,
				diff: null
			}
		},

		computed: {
			fileUpdated() {
				return (this.oldFilePresent && this.newFilePresent);
			},
			fileCreated() {
				return (!this.oldFilePresent && this.newFilePresent);
			},
			fileDeleted() {
				return (this.oldFilePresent && !this.newFilePresent);
			}
		},

		created() {
			// if this is a creation or deletion, don't display a diff

			this.oldFile = this.files.old;
			this.newFile = this.files.new;
			this.oldFilePresent = !!this.oldFile;
			this.newFilePresent = !!this.newFile;


			if (this.oldFilePresent && this.newFilePresent) {
				console.log("oldFile and newFile present, creating a diff");
				this.setupDiff(this.oldFile, this.newFile);
			};

		},

		methods: {
			setupDiff(oldFile, newFile) {
				diff = new Diff();
				let textDiff = diff.main(oldFile, newFile);
				diff.cleanupSemantic(textDiff);
				this.diff = diff.prettyHtml(textDiff);
			}
		}
	};
</script>

<style lang="scss">

	$color-updated: #001f3f;
	$color-updated-bg: lighten($color-updated, 85%);


	$color-deleted: #FF4136;
	$color-deleted-bg: lighten($color-deleted, 35%);
	$color-deleted-diff-bg: lighten($color-deleted, 25%);


	$color-created: #3D9970;
	$color-created-bg: lighten($color-created, 50%);
	$color-created-diff-bg: lighten($color-created, 25%);



	div.file {
		h2 > code {
			font-size: 86%;
		}
	}

	div {

		&.file-updated {

			h2 {
				background-color: $color-updated-bg;

				code {
					background-color: $color-updated-bg;
					color: $color-updated;
				}
			}

			pre {

				color: $color-updated;

				ins {
					background-color: $color-created-diff-bg;
					text-decoration: none;
				}

				del {
					background-color: $color-deleted-diff-bg;
					text-decoration: line-through;
				}
			}

		}

		&.file-created {

			h2 {
				background-color: $color-created-bg;

				code {
					background-color: $color-created-bg;
					color: $color-created;
				}
			}

			pre {
				color: $color-created;
			}
		}

		&.file-deleted {

			h2 {

				background-color: $color-deleted-bg;

				code {
					background-color: $color-deleted-bg;
					color: $color-deleted;
				}
			}

			pre {
				color: $color-deleted;
			}
		}
	}
</style>