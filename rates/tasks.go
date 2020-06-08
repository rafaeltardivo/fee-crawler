package rates

// Flushes cache database
func ClearCache() {
	db := newRedisDatabaseClient()
	client, err := db.getConnection()
	if err != nil {
		logger.Error(databaseError("could not connect to database"))
	}
	defer client.Close()

	_, err = client.FlushDB().Result()
	if err != nil {
		logger.Error(databaseError("failed to clear cache"))
	} else {
		logger.Info("successfully cleared cache")
	}
}
