<template>
	<div class="editor row">

		<!-- Markdown Editor Start -->
		<div class="col-md-9">
			<MarkdownEditor/>
		</div>
		<!-- Markdown Editor End -->

		<!-- Metadata Editor Start -->
		<div class="metadata-fields col-md-3">

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
		<!-- Metadata Editor End -->

	</div>

</template>

<script lang="babel">

	import MarkdownEditor from "../components/Editor/MarkdownEditor";
	import FrontMatter from "../components/Editor/FrontMatter";
	import FilenameField from "../components/Editor/FilenameField";
	import CommitMessageField from "../components/Editor/CommitMessageField";

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
