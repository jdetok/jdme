// MOST RECENT REQUESTS
db.log.aggregate([
    { $match: { url: { $regex: "^/.*players.*$" } } },
    { $sort: { request_time: -1 } },
    { $limit: 3 },
    { $project: {
        _id: 0,
        url: 1,
        time: {
            $dateToString: {
                date: "$request_time",
                format: "%m/%d/%Y %H:%M:%S",
                timezone: "America/Chicago"
            }
        }
    }
}]);
