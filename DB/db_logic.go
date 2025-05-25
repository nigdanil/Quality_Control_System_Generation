package DB

import (
	"database/sql"
	"fmt"

	"github.com/nigdanil/Quality_Control_System_Generation/model"
)

// GetLastGenerationLog загружает последнюю строку из таблицы generation_logs.
func GetLastGenerationLog(db *sql.DB) (*model.GenerationLog, error) {
	query := `
		SELECT 
			id, timestamp, 
			ref_img1_url, ref_img2_url,
			prompt, neg_prompt, seed, steps, guidance,
			ref_task1, ref_task2, ref_weight1, ref_weight2,
			width, height, ref_res,
			true_cfg, cfg_start_step, cfg_end_step, neg_guidance, first_step_guidance,
			generation_time, response_status, error_message,
			comment, effective
		FROM generation_logs
		ORDER BY id DESC
		LIMIT 1;
	`

	row := db.QueryRow(query)

	var log model.GenerationLog
	err := row.Scan(
		&log.ID, &log.Timestamp,
		&log.RefImg1URL, &log.RefImg2URL,
		&log.Prompt, &log.NegPrompt, &log.Seed, &log.Steps, &log.Guidance,
		&log.RefTask1, &log.RefTask2, &log.RefWeight1, &log.RefWeight2,
		&log.Width, &log.Height, &log.RefRes,
		&log.TrueCFG, &log.CFGStartStep, &log.CFGEndStep, &log.NegGuidance, &log.FirstStepGuidance,
		&log.GenerationTime, &log.ResponseStatus, &log.ErrorMessage,
		&log.Comment, &log.Effective,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка чтения из БД: %v", err)
	}

	return &log, nil
}
