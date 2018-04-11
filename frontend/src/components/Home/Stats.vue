<template>

	<div class="col-md-3">
		<div class="card statistics">
			<div class="card-body pb-0">
				<h4 class="card-title">Statistics</h4>

				<div class="card-block">
					<dl class="row">
						<dt class="col-sm-9">Users</dt>

						<dd class="col-sm-3 user-count">
							{{ serverInfo.users }}
						</dd>

						<dt class="col-sm-9">Commits</dt>
						<dd class="col-sm-3 commit-count">
							{{ serverInfo.commits }}
						</dd>

					</dl>
				</div>



			</div>

			<table class="table table-sm file-statistics">
				<caption class="pl-2">Files and quantities by type</caption>
				<tbody>
					<tr v-for="(value, key, i) in serverInfo.files" v-if="key != 'other'" :key="i" :data-count-type="key">
						<td class="pl-2 text-capitalize">{{ key }}</td>
						<td class="pr-2 text-center">{{ value }}</td>
					</tr>

					<tr data-count-type="other">
						<td class="pl-2 text-capitalize">Other</td>
						<td class="pr-2 text-center">{{ serverInfo.files.other }}</td>
					</tr>

				</tbody>
			</table>

		</div>
	</div>

</template>

<script>
	export default {
		name: "Stats",
		computed: {
			serverInfo() {
				return this.$store.state.server.serverInfo;
			}
		}
	}
</script>

<style lang="scss">
.file-statistics {
	td {
		padding: 0rem 1rem;
	};
};
</style>