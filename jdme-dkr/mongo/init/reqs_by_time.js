db = db.getSiblingDB("log");
// MOST RECENT REQUESTS
db.createView("rec_plr", "http", 
    [
        { $match: { url: { $regex: "^/.*players.*$" } } },
        { $sort: { request_time: -1 } },
        { $limit: 3 },
        { $project: {
                _id: 1,
                url: 1,
                time: {
                    $dateToString: {
                    date: "$request_time",
                    format: "%m/%d/%Y %H:%M:%S",
                    timezone: "America/Chicago"
                    }
                }
            }
        }
    ]
);

// db.rec_reqs.find()
db.createView(
    "rec_reqs",
    "http",
    [
        { $sort: { request_time: -1 } },
        { $limit: 3 },
        { $project: {
            _id: 1,
            url: 1,
            time: {
                $dateToString: {
                    date: "$request_time",
                    format: "%m/%d/%Y %H:%M:%S",
                    timezone: "America/Chicago"
                }
            }
        }}
    ]
);
