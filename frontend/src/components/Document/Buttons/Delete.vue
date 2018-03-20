<template>

	<div>

		<button type="button" @click="showDeleteModal" class="btn btn-danger mr-2">
			Delete
		</button>

		<div id="delete-warning" class="modal fade">
			<div class="modal-dialog" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">Delete <code>{{ document.filename }}</code></h5>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">
						<span aria-hidden="true">&times;</span>
						</button>
					</div>

					<div class="modal-body">
						Deleting a file removes it from the CMS. The deletion, along with the contents
						of the deleted files and all associated metadata <em>will be preserved</em> by
						the underlying version control system.
					</div>

					<div class="modal-body">
						<p>
							Are you sure you want to delete <code>{{ document.filename }}</code>?
						</p>

						<div class="form-check">
							<label class="form-check-label">
								<input v-model="deleteAttachments" class="form-check-input" type="checkbox">
								Delete attachments
							</label>

							<p class="form-text text-muted">
								By default attachments are left in the repository when a file is deleted.
							</p>
						</div>

					</div>
					<div class="modal-footer">

						<button type="button" @click="remove" class="btn btn-danger mr-2">
							Confirm deletion
						</button>

						<button class="btn btn-secondary" data-dismiss="modal">
							Cancel
						</button>
					</div>
				</div>
			</div>
		</div>

	</div>

</template>

<script lang="babel">

	import Accessors from '../../Mixins/accessors';

	import CMSDirectory from '../../../javascripts/models/directory.js';
	import checkResponse from "../../../javascripts/response.js";

	export default {
		name: "DocumentDelete",
		data() {
			return {
				deleteAttachments: false
			};
		},
		created() {
			// create a commit to be populated/used if delete is clicked
			this.initializeCommit();
		},
		methods: {
			initializeCommit() {
				this.$store.dispatch("initializeCommit");
			},

			showDeleteModal() {
				event.preventDefault();
				return $("#delete-warning.modal").modal();
			},

			hideDeleteModal() {
				return $("#delete-warning.modal").modal("hide");
			},

			async remove(event) {

				event.preventDefault();

				// prepare the commit
				this.commit.addFile(this.document);
				if (this.deleteAttachments) {
					this.commit.addDirectory(new CMSDirectory(this.document.attachmentsDir));
				};

				// try to destroy
				let response = await this.destroy();

				// if the repo has been updated in the background, refresh the
				// document and try again
				if (response.status == 409) {
					await this.getDocument();
					response = await this.destroy();
				};

				if (!checkResponse(response.status)) {
					console.error("Could not delete document", response);
					return;
				};

				// deletion was successful, update the latest rev and redirect
				let json = await response.json();
				this.$store.commit("setLatestRevision", json.oid);
				this.hideDeleteModal();
				this.redirectToDirectoryIndex(this.params.directory);

			},

			async destroy() {
				let response = await this.document.destroy(this.commit);
				return response;
			},

			redirectToDirectoryIndex(directory) {
				this.$router.push({
					name: 'directory_index',
					params:{directory}
				});
			},
			async getDocument() {
				let filename = "index.md";

				if (this.params.language_code) {
					filename = `index.${this.params.language_code}.md`;
				};

				let directory = this.params.directory;
				let document = this.params.document;

				return this.$store.dispatch("getDocument", {directory, document, filename});
			},
		},
		mixins: [Accessors]
	};
</script>