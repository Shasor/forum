package models

import (
	"login/src/database"
	"strings"
)

func CountLikeDislikePost(post database.Post) (int, int) {
	TabLike := strings.Split(post.Like, "-")
	TabDislike := strings.Split(post.Dislike, "-")
	return len(TabLike) - 1, len(TabDislike) - 1
}

func CountLikeDislikeCommentaire(commentaire database.Commentaire) (int, int) {
	TabLike := strings.Split(commentaire.Like, "-")
	TabDislike := strings.Split(commentaire.Dislike, "-")
	return len(TabLike) - 1, len(TabDislike) - 1
}
