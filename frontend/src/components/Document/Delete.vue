<template>

	<div>

		<!-- FIXME had to remove 'fade' class from modal due to Cucumber, add it back -->
		<div id="delete-warning" class="modal">
			<div class="modal-dialog" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">Delete <code>{{ document.filename }}</code></h5>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">
						<span aria-hidden="true">&times;</span>
						</button>
					</div>

					<div class="modal-body">
						Deleting a file removes it from the CMS, but don't worry, an administrator or
						technical user can restore files in the future.
					</div>

					<div class="modal-body" v-if="attachmentsPresent">
						<p>
							Are you sure you want to delete <code>{{ document.filename }}</code>?
						</p>

						<p class="text-muted">
							By default attachments are left in the repository when a file is deleted.
							If you want to delete them, please check the box below.
						</p>

						<div class="form-check">
							<label class="form-check-label">
								<input v-model="deleteAttachments" class="form-check-input" type="checkbox">
								Delete attachments?
							</label>
						</div>

					</div>
					<div class="modal-footer">

						<button type="button" @click="destroy" class="btn btn-danger mr-2">
							Confirm deletion
						</button>

						<button class="btn btn-secondary" data-dismiss="modal">
							Cancel
						</button>
					</div>
				</div>
			</div>
		</div>


		<button type="button" @click="showDeleteModal" class="btn btn-danger mr-2">
			Delete
		</button>

	</div>

</template>

<script lang="babel">

	import Accessors from '../Mixins/accessors';

	import CMSDirectory from '../../javascripts/models/directory.js';
	import checkResponse from "../../javascripts/response.js";

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
				$("#delete-warning.modal").modal();
			},

			hideDeleteModal() {
				$("#delete-warning.modal").modal("hide");
			},

			async destroy(event, ) {

				event.preventDefault();

				// build and configure the commit 👷‍

				this.commit.addFile(this.document);

				if (this.deleteAttachments) {
					this.commit.addDirectory(new CMSDirectory(this.document.attachmentsDir));
				};

				let response = await this.document.destroy(this.commit, false);

				if (!checkResponse(response.status)) {

					if (response.status == 409) {

						this.$store.state.broadcast.addMessage(
							"danger",
							"Failed",
							"The repository is out of sync",
							3
						);

						this.hideDeleteModal();
						this.getDocument();
						this.commit.reset();

						return;
					};

					// any other error
					throw("could not delete document", response);
					return;
				};

				this.hideDeleteModal();
				this.redirectToDirectoryIndex(this.directory);

			},
			redirectToDirectoryIndex(directory) {
				this.$router.push({
					name: 'document_index',
					params:{directory}
				});
			},
			async getDocument() {
				let filename = this.filename;
				let directory = this.directory;
				this.$store.dispatch("getDocument", {directory, filename});
			},
		},
		computed: {
			attachmentsCount() {
				return this.document.attachments && this.document.attachments.length;
			},
			attachmentsPresent() {
				return this.attachmentsCount && this.attachmentsCount > 0;
			}
		},
		mixins: [Accessors]
	};
</script>