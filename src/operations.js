'use strict';

const assert = require('assert');

const { hasNeedRebaseLabel, LABEL_STATUS_NEED_REBASE } = require('./labels');

async function getOpenPullRequestAll(octokit, owner, repo) {
    // FIXME: this code could not get all opened pull requests if there are over 100.
    const result = await octokit.pulls.list({
        owner,
        repo,
        state: 'open',
        // eslint-disable-next-line camelcase
        per_page: '100',
    });
    const body = result.data;
    assert.strictEqual(!!body, true);

    return body;
}

function getPullRequestId(pullReqInfo) {
    const id = pullReqInfo.number;
    assertTypeof(id, 'number');
    return id;
}

const SECONDS_WE_ASSUME_GITHUB_WOULD_COMPLETE_CHECK_IF_MERGEBLE = 1000 * 10;

async function checkAndMarkIfPullRequestUnmergeable(octokit, owner, repo, oldPullReqInfo, compareUrl) {
    const number = getPullRequestId(oldPullReqInfo);
    if (hasNeedRebaseLabel(oldPullReqInfo.labels)) {
        console.log(`#${number} has been labeled as ${LABEL_STATUS_NEED_REBASE}`);
        return;
    }

    // When we get all opened pull requests, GitHub have not checked whether ths PR is unmergeble yet.
    // So I think we should retry them.
    let mergeable = await shouldMarkPullRequestNeedRebase(octokit, owner, repo, number);
    if (mergeable === null) {
        // retry
        await sleep(SECONDS_WE_ASSUME_GITHUB_WOULD_COMPLETE_CHECK_IF_MERGEBLE);
        mergeable = await shouldMarkPullRequestNeedRebase(octokit, owner, repo, number);
    }

    console.log(`The mergeable of #${number} is ${mergeable}`);
    if (mergeable === null) {
        // give up
        return;
    }

    if (mergeable) {
        return;
    }

    const addLabel = octokit.issues.addLabels({
        owner,
        repo,
        number,
        labels: [LABEL_STATUS_NEED_REBASE],
    });

    const body = `:umbrella: The latest upstream change (presumably [these](${compareUrl})) made this pull request unmergeable. Please resolve the merge conflicts.`;
    const addComment = octokit.issues.createComment({
        owner,
        repo,
        number,
        body,
    });

    await Promise.all([addLabel, addComment]);
}

async function shouldMarkPullRequestNeedRebase(octokit, owner, repo, number) {
    const newPRInfoResponse = await octokit.pulls.get({
        owner,
        repo,
        number,
    });
    const newPRInfo = newPRInfoResponse.data;
    assert.strictEqual(!!newPRInfo, true);

    // Check again to confirm the other instance of this action's behavior.
    if (hasNeedRebaseLabel(newPRInfo.labels)) {
        return false;
    }

    const mergeable = newPRInfo.mergeable;
    if (mergeable === null) {
        return null;
    }

    return mergeable;
}

function sleep(millisec) {
    const p = new Promise((resolve) => {
        setTimeout(resolve, millisec);
    });
    return p;
}

function assertTypeof(val, typename, message) {
    assert.strictEqual(typeof val, typename, message);
}

function isRelatedToPushedBranch(pullReqInfo, pushedBranchRef) {
    //  Define to explain followings comments:
    //      * _origin_ is the repository which registers & runs this action.
    //      * _forked_ is the repository which is forked from _origin_.
    //
    //  We don't have to imagine the case which the user pushed to the branch in _forked_.
    //  and its branch is opened as the pull requet for _origin_.
    //  Because _push_ event happens and this action will run only if someone pushed to a branch on _origin_.
    //
    //  By the observation of the behavior of GitHub REST API v3 Pull Requests,
    //  object's values would be followings if `repository_name:branch_name` is `octocat:master`.
    //
    //      * `.base.ref`
    //          * The value is `master`.
    //          * This value is the branch name which the pull request would be merged into.
    //            even if the pull request is came from _forked_ or the one's parent is arbitary commit.
    //      * `.base.label` is `octocat:master`.
    //          * The value is `octocat:master`.
    //          * This value is same even if the pull request is came from _forked_ or the one's parent is arbitary commit.
    //
    //  By theese result, I think we can judge that the pull request is related to the pushed branch by checking
    //  whether `.base.ref` is same with the name of the pushed branch.

    const currentBranchName = pullReqInfo.base.ref; // e.g. `master` in `octocat:master`
    const pushedBranchName = pushedBranchRef.replace('refs/heads/', '');
    const hasRelationShip = currentBranchName === pushedBranchName;

    return hasRelationShip;
}

module.exports = Object.freeze({
    getOpenPullRequestAll,
    checkAndMarkIfPullRequestUnmergeable,
    isRelatedToPushedBranch,
    getPullRequestId,
});
