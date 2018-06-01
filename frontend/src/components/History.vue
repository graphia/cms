<template>
	<div id="repo-history">

		<Breadcrumbs :levels="breadcrumbs"/>

		<div class="col col-md-12">
			<h1>History</h1>

			<div class="commit-list">
				<div class="card m-4" v-for="(commit, i) in commits" :key="i">
					<div class="card-header">
						<div class="hash">
							{{ prettyDate(commit.timestamp) }} <span class="text-muted">({{ timeAgo(commit.timestamp) }})</span>
						</div>
						<div class="author">
							{{ commit.author.Name }}
						</div>
					</div>
					<div class="card-body">

						<router-link class="btn btn-text" :to="{name: 'commit', params: {hash: commit.id}}">
							{{ commit.message }}
						</router-link>

					</div>
				</div>
			</div>
		</div>
	</div>
</template>


<script lang="babel">

	// external
	import fecha from 'fecha';
	import vagueTime from 'vague-time';

	// javascripts
	import config from '../javascripts/config.js';
	import checkResponse from '../javascripts/response.js';
	import CMSBreadcrumb from '../javascripts/models/breadcrumb.js';

	// components
	import Breadcrumbs from './Utilities/Breadcrumbs';


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
			},
			parseDate(d) {
				return fecha.parse(d, 'YYYY-MM-DDTHH:mm');
			},
			prettyDate(d) {
				let value = this.parseDate(d);

				return fecha.format(value, 'YYYY-MM-DD HH:mm');
			},
			timeAgo(d) {
				let value = this.parseDate(d);

				return vagueTime.get({
					from: Date.now(),
					to: Date.parse(value),
					units: 'ms'
				});

			}
		},
		computed: {
			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						"History",
						"history"
					)
				];
			},
		},
		components: {
			Breadcrumbs
		}
	};
</script>

<style lang="scss" scoped>

	.commit-list {

		border-left: 3px solid grey;

		.card > .card-header {
			display: flex;

			.hash {
				flex-grow: 1;
			};

		}
	}
</style>