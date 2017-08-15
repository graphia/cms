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
			<CommitMessageField/>

			<div class="form-group">

				<input
					type="submit"
					value="Update"
					class="btn btn-success"
					v-bind:disabled="!valid"
				/>

				<router-link :to="{name: 'document_show', params: {directory: 'documents', filename: document.filename}}" class="btn btn-text">
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
		components: {
			MarkdownEditor,
			FrontMatter,
			CommitMessageField,
		},
		methods: {
			validate() {
				if (!this.form) {
					this.form = document.getElementById("edit-document-form");
				};
				this.valid = this.form.checkValidity();;
			}
		}
	};
</script>
