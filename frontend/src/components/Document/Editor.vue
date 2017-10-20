<template>
	<div class="editor row">

		<!-- Markdown Editor Start -->
		<div class="col-md-8">
			<MarkdownEditor/>
		</div>
		<!-- Markdown Editor End -->

		<!-- Sidebar Start -->
		<div class="col-md-4">
			<div class="sidebar card">

				<div class="card-header">

					<ul class="nav nav-tabs card-header-tabs" role="tablist">
						<li class="nav-item">
							<a class="nav-link active" role="tab" data-toggle="tab" href="#metadata">Document Info</a>
						</li>

						<li class="nav-item">
							<a class="nav-link" role="tab" data-toggle="tab" href="#gallery">Images</a>
						</li>
					</ul>

				</div>

				<div class="tab-content">

					<div id="metadata" class="active tab-pane card-body metadata-fields" role="tab-panel">

						<FrontMatter/>

						<FilenameField v-if="newFile"/>

						<CommitMessageField/>

						<div class="form-group">

							<input
								type="submit"
								class="btn btn-success"
								:value="submitButtonText"
								v-bind:disabled="!valid"
							/>

							<router-link :to="formCancellationRedirectParams" class="btn btn-text">
								Cancel
							</router-link>
						</div>
					</div>

					<Gallery id="gallery" class="tab-pane card-body" role="tab-panel"/>

				</div>
			</div>
			<!-- Sidebar Editor End -->
		</div>

	</div>

</template>

<script lang="babel">

	import MarkdownEditor from "./Editor/MarkdownEditor";
	import FrontMatter from "./Editor/FrontMatter";
	import FilenameField from "./Editor/FilenameField";
	import Gallery from "./Editor/Gallery";
	import CommitMessageField from "./Editor/CommitMessageField";

	export default {
		name: "Editor",

		created() {
			this.$bus.$on("checkMetadata", () => {
				this.validate()
			});
		},

		data() {
			return {
				markdownLoaded: false,
				valid: false
			};
		},
		computed: {
			// quick access to things in the store
			document() {
				return this.$store.state.activeDocument;
			}
		},
		watch: {
			// FIXME use bus instead of cascading
			"$parent.markdownLoaded": function() {
				this.markdownLoaded = true;
			}
		},
		props: [
			'formID',
			'submitButtonText',
			'newFile',
			'formCancellationRedirectParams'
		],
		components: {
			MarkdownEditor,
			FrontMatter,
			FilenameField,
			CommitMessageField,
			Gallery
		},
		methods: {
			validate() {
				if (!this.form) {
					this.form = document.getElementById(this.formID);
				};
				this.valid = this.form.checkValidity();
			}
		}
	};
</script>
