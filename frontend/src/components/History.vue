<template>
	<div class="row">

		<div class="col col-md-12">
			<h1>History</h1>

			<div class="commit-list">
				<div class="card m-4" v-for="(commit, i) in commits" :key="i">
					<div class="card-header">
						<div class="hash">
							{{ commit.id }}
						</div>
						<div class="author">
							{{ commit.author.Name }}
						</div>
					</div>
					<div class="card-body">
						{{ commit.message }}

					</div>
				</div>
			</div>
		</div>
	</div>
</template>


<script lang="babel">

	import config from '../javascripts/config.js';
	import checkResponse from '../javascripts/response.js';

	export default {
		name: "History",
		created() {
			this.getHistory();
		},
		data() {
			return {
				commits: []
			};
		},
		methods: {
			async getHistory() {
				const path = `${config.api}/history`;
				let response = await fetch(path, {headers: this.$store.state.auth.authHeader()})

				if (!checkResponse(response.status)) {
						console.error("History cannot be retrieved", response);
						return;
				};

				this.commits = await response.json();
			}
		}
	};
</script>

<style lang="scss" scoped>

	.commit-list {

		border-left: 3px solid red;

		.card > .card-header {
			display: flex;

			.hash {
				flex-grow: 1;
			};

		}
	}
</style>