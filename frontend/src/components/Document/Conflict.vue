<script lang="babel">
	export default {
		name: "DocumentConflict",
	};
</script>

<template>

	<div id="conflict-warning" class="modal fade">
		<div class="modal-dialog" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Your commit failed</h5>
					<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
					</button>
				</div>
				<div class="modal-body">
					<p>
						The repository has been modified since you began editing.
					</p>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-primary">Download your copy</button>
					<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
				</div>
			</div>
		</div>
	</div>

</template>