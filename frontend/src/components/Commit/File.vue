<template>
	<div class="card file" v-bind:class="{
			'file-updated card-outline-info': fileUpdated,
			'file-created card-outline-success': fileCreated,
			'file-deleted card-outline-danger': fileDeleted
		}"
	>

		<h2 class="card-header">
			<octicon :icon-name="this.icon">omg</octicon>
			<code>{{ this.path }}</code>
		</h2>

		<div class="card-block">

			<!-- file content display, different content shown depending on the git operation -->

			<div v-if="this.fileUpdated">
				<pre v-html="this.diff"/>
			</div>

			<div v-else-if="this.fileCreated" class="file-created">
				<pre>{{ this.newFile }}</pre>
			</div>

			<div v-else-if="this.fileDeleted" class="file-deleted">
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
			},
			icon() {
				let text = null;

				switch (true) {
					case this.fileUpdated:
						text = "diff-modified";
						break;
					case this.fileCreated:
						text = "diff-added";
						break;
					case this.fileDeleted:
						text = "diff-removed";
						break;
				};

				return text;
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

	// updates/modifications
	$color-updated: #001f3f;
	$color-updated-bg: lighten($color-updated, 85%);

	// deletions
	$color-deleted: #FF4136;
	$color-deleted-bg: lighten($color-deleted, 35%);
	$color-deleted-diff-bg: lighten($color-deleted, 35%);

	// creation
	$color-created: #3D9970;
	$color-created-bg: lighten($color-created, 50%);
	$color-created-diff-bg: lighten($color-created, 40%);

	div.file {

		h2 > code {
			font-size: 70%;
			color: inherit !important;
		}

		pre {
			white-space: pre-wrap !important;
		}
	}

	div {

		&.file-updated {

			h2 {
				background-color: $color-updated-bg;
				color: $color-updated;

				code {
					background-color: $color-updated-bg;
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
				color: $color-created;

				code {
					background-color: $color-created-bg;

				}
			}

			pre {
				color: $color-created;
			}
		}

		&.file-deleted {

			h2 {

				background-color: $color-deleted-bg;
				color: $color-deleted;

				code {
					background-color: $color-deleted-bg;

				}
			}

			pre {
				color: $color-deleted;
			}
		}
	}
</style>