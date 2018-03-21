<template>

	<span>

		<button @click="showDeleteModal" class="btn btn-danger btn-sm">
			Delete directory
		</button>

		<div id="delete-warning" class="modal">
			<div class="modal-dialog" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">Delete {{ activeDirectory.title }}</h5>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">
						<span aria-hidden="true">&times;</span>
						</button>
					</div>

					<div class="modal-body">
						<p>
							Are you sure you want to delete {{ activeDirectory.title }} and all
							of its contents?
						</p>

					</div>
					<div class="modal-footer">

						<button type="button" @click="deleteDirectory" class="btn btn-danger mr-2">
							Confirm deletion
						</button>

						<button class="btn btn-secondary" data-dismiss="modal">
							Cancel
						</button>
					</div>
				</div>
			</div>
		</div>

	</span>

</template>


<script lang="babel">

	import config from '../../../javascripts/config.js';
	import CMSCommit from '../../../javascripts/models/commit.js';
	import checkResponse from '../../../javascripts/response.js';

	import Accessors from '../../Mixins/accessors';

	export default {
		name: "DirectoryDelete",
		mixins: [Accessors],
		methods: {


			showDeleteModal(event) {
				event.preventDefault();
				return $("#delete-warning.modal").modal();
			},

			hideDeleteModal() {
				return $("#delete-warning.modal").modal("hide");
			},

			async deleteDirectory(event) {
				event.preventDefault();
				this.hideDeleteModal();

				let title = this.activeDirectory.title

				let commit = new CMSCommit(
					`Deleting directory ${title}`,
					[], // no files
					[this.activeDirectory]
				);

				let response = await this.activeDirectory.destroy(commit);

				if (!checkResponse(response.status)) {
					console.error("could not delete directory", response);
					return;
				};

				// success
				let json = await response.json();

				await this.$store.commit("setLatestRevision", json.oid);

				this.$store.state.broadcast.addMessage(
					"success",
					"Directory deleted",
					`${title} and its contents have been deleted`,
					3
				);

				this.redirectToHome();

			},

			redirectToHome() {
				this.$router.push({name: 'home'});
			}
		}
	};
</script>
