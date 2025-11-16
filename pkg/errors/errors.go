package errors

type BadRequestError struct{}

func (e *BadRequestError) Error() string {
	return "BAD_REQUEST"
}

type TeamExistsError struct{}

func (e *TeamExistsError) Error() string {
	return "team_name already exists"
}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "NOT_FOUND"
}

type PrExistsError struct{}

func (e *PrExistsError) Error() string {
	return "PR id already exists"
}

type NoActiveReviewersError struct{}

func (e *NoActiveReviewersError) Error() string {
	return "NO_ACTIVE_REVIEWERS"
}

type InternalServerError struct{}

func (e *InternalServerError) Error() string {
	return "INTERVAL_SERVER_ERROR"
}

type PrMergedError struct{}

func (e *PrMergedError) Error() string {
	return "cannot reassign on merged PR"
}

type NotAssignedError struct{}

func (e *NotAssignedError) Error() string {
	return "reviewer is not assigned to this PR"
}

type NoCandidateError struct{}

func (e *NoCandidateError) Error() string {
	return "no active replacement candidate in team"
}
