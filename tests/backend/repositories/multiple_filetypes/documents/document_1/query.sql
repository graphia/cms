select *
from reporters
where name not in (
	"Kent Brockman",
	"Arnie Pye",
	"Dave Shutton",
	"Chloe Talbot"
);
